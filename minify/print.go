package minify

import (
	"fmt"
	"io"
	"reflect"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

func Print(w io.Writer, m *javascript.Module) (int64, error) {
	var h walk.Handler

	h = walk.HandlerFunc(func(t javascript.Type) error {
		switch t := t.(type) {
		case *javascript.PropertyDefinition:
			if !t.IsCoverInitializedName && t.PropertyName != nil && t.PropertyName.LiteralPropertyName != nil && t.AssignmentExpression != nil && t.AssignmentExpression.ConditionalExpression != nil {
				c := javascript.UnwrapConditional(t.AssignmentExpression.ConditionalExpression)

				if pe, ok := c.(*javascript.PrimaryExpression); ok && pe.IdentifierReference != nil && pe.IdentifierReference.Type == t.PropertyName.LiteralPropertyName.Type && pe.IdentifierReference.Data == t.PropertyName.LiteralPropertyName.Data {
					v := *t.PropertyName.LiteralPropertyName
					pe.IdentifierReference = &v
				}
			}
		case *javascript.AssignmentProperty:
			if t.PropertyName.LiteralPropertyName != nil && t.DestructuringAssignmentTarget != nil && t.DestructuringAssignmentTarget.LeftHandSideExpression != nil && t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression != nil && len(t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News) == 0 && t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression != nil && t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil && t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference.Data == t.PropertyName.LiteralPropertyName.Data {
				v := *t.PropertyName.LiteralPropertyName
				t.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference = &v
			}
		case *javascript.BindingProperty:
			if t.PropertyName.LiteralPropertyName != nil && t.BindingElement.SingleNameBinding != nil && t.PropertyName.LiteralPropertyName.Data == t.BindingElement.SingleNameBinding.Data {
				v := *t.PropertyName.LiteralPropertyName
				t.BindingElement.SingleNameBinding = &v
			}
		case *javascript.ArrayAssignmentPattern:
			if t.AssignmentRestElement != nil {
				break
			}

			for n := len(t.AssignmentElements) - 1; n >= 0; n-- {
				if t.AssignmentElements[n].DestructuringAssignmentTarget.AssignmentPattern == nil && t.AssignmentElements[n].DestructuringAssignmentTarget.LeftHandSideExpression == nil {
					t.AssignmentElements = t.AssignmentElements[:n]
				} else {
					break
				}
			}
		case *javascript.ArrayBindingPattern:
			if t.BindingRestElement != nil {
				break
			}

			for n := len(t.BindingElementList) - 1; n >= 0; n-- {
				if (t.BindingElementList[n].SingleNameBinding == nil || t.BindingElementList[n].SingleNameBinding.Data == "") && t.BindingElementList[n].ArrayBindingPattern == nil && t.BindingElementList[n].ObjectBindingPattern == nil {
					t.BindingElementList = t.BindingElementList[:n]
				} else {
					break
				}
			}
		}

		v := reflect.ValueOf(t)

		if v.Type().Kind() != reflect.Pointer || v.Type().Elem().Kind() != reflect.Struct {
			return nil
		}

		if f, ok := v.Type().Elem().FieldByName("Tokens"); ok {
			v.Elem().FieldByIndex(f.Index).SetZero()
		}

		return walk.Walk(t, h)
	})

	h.Handle(m)

	n, err := fmt.Fprintf(w, "%#s", m)

	return int64(n), err
}
