use crate::snippet::snippet;

macro_rules! expression_snippet_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let (expression, pos, expected_snippet) = $value;
                let mut obtained_snippet = String::from('\n');
            	obtained_snippet.push_str(snippet(expression.to_string(), pos).as_str());
            	assert_eq!(expected_snippet.to_string(), obtained_snippet);
			}
		)*
		}
	}

expression_snippet_tests! {
    expr_0: ("012345678", 5, r#"
| 012345678
| .....^"#),
    expr_1: ("012345678", 0, r#"
| 012345678
| ^"#),
    expr_2: ("012345678", 10, r#"
| 012345678
| ........^"#),
    expr_3: ("", 10, "\n"),
}
