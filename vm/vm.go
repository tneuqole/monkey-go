package vm

import (
	"fmt"

	"github.com/tneuqole/monkey-go/code"
	"github.com/tneuqole/monkey-go/compiler"
	"github.com/tneuqole/monkey-go/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	globals      []object.Object

	// always points to next free space
	// top of stack is stack[sp-1]
	sp int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		globals:      make([]object.Object, GlobalsSize),
		sp:           0,
	}
}

func NewWithGlobals(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		var err error
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err = vm.push(vm.constants[constIndex])
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err = vm.executeBinaryOperation(op)
		case code.OpTrue:
			err = vm.push(True)
		case code.OpFalse:
			err = vm.push(False)
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err = vm.executeComparison(op)
		case code.OpBang:
			err = vm.executeBangOperator()
		case code.OpMinus:
			err = vm.executeMinusOperator()
		case code.OpPop:
			vm.pop()
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			// -1 because ip is incremented after the loop
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err = vm.push(vm.globals[globalIndex])
		case code.OpNull:
			err = vm.push(Null)
		}

		if err != nil {
			return err
		}

	}

	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("STACK OVERFLOW")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	val := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -val})
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	} else if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	var result int64
	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal
	default:
		return fmt.Errorf("unknown integer operater: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	var result string
	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	default:
		return fmt.Errorf("unknown string operater: %d", op)
	}

	return vm.push(&object.String{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal == rightVal))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal != rightVal))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown integer operater: %d", op)
	}
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
