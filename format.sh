#!/bin/bash

types() {
	for file in ast_class.go ast_conditional.go ast_expression.go ast_function.go ast.go ast_module.go ast_statement.go ; do
		while read type; do
			echo "$type" "$file";
		done < <(grep "type .* struct {" "$file" | cut -d' ' -f2);
	done | sort;
}

(
	cat <<HEREDOC
package javascript

// File automatically generated with format.sh.

import "io"
HEREDOC

	while read type file; do
		echo -e "\nfunc (f *$type) printType(w io.Writer, v bool) {";
		echo "	pp := indentPrinter{w}";
		echo;
		echo "	pp.Print(\"$type {\")";
		while read fieldName fieldType; do
			if [ -z "$fieldType" ]; then
				fieldType="$fieldName";
				fieldName="${fieldName/\*/}";
			fi;

			if [ "$fieldType" = "bool" ]; then
				echo;
				echo "	if f.$fieldName || v {";
				echo "		pp.Printf(\"\\n$fieldName: %v\", f.$fieldName)";
				echo "	}";
			elif [ "$fieldType" = "uint" -o "$fieldType" = "int" ]; then
				echo;
				echo "	if f.$fieldName != 0 || v {";
				echo "		pp.Printf(\"\\n$fieldName: %v\", f.$fieldName)";
				echo "	}";
			elif [ "${fieldType:0:2}" = "[]" ]; then
				echo;
				echo "	if f.$fieldName == nil {";
				echo "		pp.Print(\"\\n$fieldName: nil\")";
				echo "	} else if len(f.$fieldName) > 0 {";
				echo "		pp.Print(\"\\n$fieldName: [\")";
				echo;
				echo "		ipp := indentPrinter{&pp}";
				echo;
				echo "		for n, e := range f.$fieldName {";
				echo "			ipp.Printf(\"\n%d: \", n)";
				echo "			e.printType(&ipp, v)";
				echo "		}";
				echo;
				echo "		pp.Print(\"\\n]\")";
				echo "	} else if v {";
				echo "		pp.Print(\"\\n$fieldName: []\")";
				echo "	}";
			elif [ "${fieldType:0:1}" = "*" -a "$fieldType" != "*Token" ]; then
				echo;
				echo "	if f.$fieldName != nil {";
				echo "		pp.Print(\"\\n$fieldName: \")";
				echo "		f.$fieldName.printType(&pp, v)";
				echo "	} else if v {";
				echo "		pp.Print(\"\\n$fieldName: nil\")";
				echo "	}";
			else
				echo;
				echo "	pp.Print(\"\\n$fieldName: \")";
				echo "	f.$fieldName.printType(&pp, v)";
			fi;
		done < <(sed '/^type '$type' struct {$/,/^}$/!d;//d' "$file" | grep -v "^$");

		echo;
		echo "	io.WriteString(w, \"}\")";
		echo "}";
	done < <(types);
) > "format_types.go";
