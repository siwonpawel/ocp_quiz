package main

import (
	"reflect"
	"testing"
)

func Test_addCustomClasses(t *testing.T) {
	type args struct {
		value []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Should do nothing to empty value",
			args: args{value: []byte("")},
			want: []byte(""),
		},
		{
			name: "Should add custom classes to `class Person {`",
			args: args{value: []byte("class Person")},
			want: []byte(`<p class="cb-type">class</p> <p class="cb-type-name">Person</p>`),
		},
		{
			name: "Should add custom classes to `private class Person`",
			args: args{value: []byte("private class Person")},
			want: []byte(`<p class="cb-access-modifier">private</p> <p class="cb-type">class</p> <p class="cb-type-name">Person</p>`),
		},
		{
			name: "Should add custom classes to `private final class Person`",
			args: args{value: []byte("private final class Person")},
			want: []byte(`<p class="cb-access-modifier">private</p> <p class="cb-final">final</p> <p class="cb-type">class</p> <p class="cb-type-name">Person</p>`),
		},
		{
			name: "Should add custom classes to `public abstract class Person`",
			args: args{value: []byte("public abstract class Person")},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-abstract">abstract</p> <p class="cb-type">class</p> <p class="cb-type-name">Person</p>`),
		},
		{
			name: "Should add custom classes to `interface Serializable`",
			args: args{value: []byte("interface Serializable")},
			want: []byte(`<p class="cb-type">interface</p> <p class="cb-type-name">Serializable</p>`),
		},
		{
			name: "Should add custom classes to `public interface Serializable`",
			args: args{value: []byte("public interface Serializable")},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-type">interface</p> <p class="cb-type-name">Serializable</p>`),
		},
		{
			name: "Should add custom classes to `final public interface Serializable`",
			args: args{value: []byte("final public interface Serializable")},
			want: []byte(`<p class="cb-final">final</p> <p class="cb-access-modifier">public</p> <p class="cb-type">interface</p> <p class="cb-type-name">Serializable</p>`),
		},
		{
			name: "Should add custom classes to `enum TimeUnit`",
			args: args{value: []byte("enum TimeUnit")},
			want: []byte(`<p class="cb-type">enum</p> <p class="cb-type-name">TimeUnit</p>`),
		},
		{
			name: "Should add custom classes to `public enum TimeUnit`",
			args: args{value: []byte("public enum TimeUnit")},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-type">enum</p> <p class="cb-type-name">TimeUnit</p>`),
		},
		{
			name: "Should add custom classes to `record Animal`",
			args: args{value: []byte("record Animal")},
			want: []byte(`<p class="cb-type">record</p> <p class="cb-type-name">Animal</p>`),
		},
		{
			name: "Should add custom classes to `public record Animal`",
			args: args{value: []byte("public record Animal")},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-type">record</p> <p class="cb-type-name">Animal</p>`),
		},
		{
			name: "Should add custom classes to `public class Person {\n\tprivate final String name;\n}",
			args: args{value: []byte(`public class Person {\\\n\tprivate final String name;\\\n}`)},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-type">class</p> <p class="cb-type-name">Person</p> {\n\t<p class="cb-access-modifier">private</p> <p class="cb-final">final</p> <p class="cb-type-name">String</p> name;\n}`),
		},
		{
			name: "Should add custom classes to 'public class Person {\n\tprivate final String name;\n\tpublic String getName() {\n\t\treturn name;\n\t}\n}`",
			args: args{value: []byte(`public class Person {\\n\\tprivate final String name;\\n\\tpublic String getName() {\\n\\t\\treturn name;\\n\\t}\\n}`)},
			want: []byte(`<p class="cb-access-modifier">public</p> <p class="cb-type">class</p> <p class="cb-type-name">Person</p> {\n\t<p class="cb-access-modifier">private</p> <p class="cb-final">final</p> <p class="cb-type-name">String</p> name;\n\t<p class="cb-access-modifier">public</p> <p class="cb-type-name">String</p> <p class="cb-method">getName()</p> {\n\t\treturn name;\n\t}\n}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addCustomClasses(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addCustomClasses() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
