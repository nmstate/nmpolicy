use crate::{
    error::{ErrorKind, NmpolicyError},
    lex::{
        tokens::Token,
        tokens::Token::{Dot, Identity, Merge, Number, Pipe, Replace, Str},
        tokens::Tokens,
    },
};

type Results<'a> = Vec<<Tokens<'a> as Iterator>::Item>;

macro_rules! basic_expression_tests{
		($($name:ident: $value:expr,)*) => {
		$(
			#[test]
			fn $name() {
                let (expression, expected_tokens) = $value;
            	println!("{}", expression);
            	let tokens = Tokens::new(expression);
                let obtained_result: Results = tokens.collect();
                let expected_result: Results = expected_tokens.iter().map(|t| Ok(t.clone())).collect();
            	assert_eq!(expected_result, obtained_result);
			}
		)*
		}
	}
basic_expression_tests! {
    expr_0: ("", Vec::<Token>::new()),
    expr_1: ("    ", Vec::<Token>::new()),
    expr_2: ("    31    03   ", vec![
        Number(4, 31),
        Number(10, 3),
    ]),
    expr_3: (r#" "foobar1" "foo 1 bar"    " foo bar - " ' bar foo' 789 "" "#, vec![
        Str(1, "foobar1".to_string()),
        Str(11, "foo 1 bar".to_string()),
        Str(26, " foo bar - ".to_string()),
        Str(40, " bar foo".to_string()),
        Number(51, 789),
        Str(55, "".to_string()),
    ]),
    expr_4: (r#" . foo1.dar1.0.dar2:=foo3 . dar3 ... moo3+boo3|doo3"#, vec![
        Dot(1),
        Identity(3, "foo1".to_string()),
        Dot(7),
        Identity(8, "dar1".to_string()),
        Dot(12),
        Number(13, 0),
        Dot(14),
        Identity(15, "dar2".to_string()),
        Replace(19),
        Identity(21, "foo3".to_string()),
        Dot(26),
        Identity(28, "dar3".to_string()),
        Dot(33),
        Dot(34),
        Dot(35),
        Identity(37, "moo3".to_string()),
        Merge(41),
        Identity(42, "boo3".to_string()),
        Pipe(46),
        Identity(47, "doo3".to_string()),
    ]),
}
macro_rules! failure_tests{
        ($($name:ident: $value:expr,)*) => {
        $(
            #[test]
            fn $name() {
                let (expression, expected_error) = $value;
                println!("{}", expression);
                let mut tokens = Tokens::new(expression);
                let obtained_error = tokens.find(Result::is_err);
                assert_eq!(Some(Err(NmpolicyError::new(ErrorKind::InvalidArgument, expected_error))), obtained_error);
                assert_eq!(None, tokens.next());
            }
        )*
        }
    }
failure_tests! {
    err_1: ("foo=bar", r#"invalid EQFILTER operation format (b is not equal char)
| foo=bar
| ....^"#.to_string()),
    err_2: (" foo 1foo ", r#"invalid number format (f is not a digit)
|  foo 1foo 
| ......^"#.to_string()),
    err_3: (" foo -foo ", r#"illegal char -
|  foo -foo 
| .....^"#.to_string()),
    err_4: (r#" "bar1" "foo dar"#, r#"invalid string format (missing " terminator)
|  "bar1" "foo dar
| ...............^"#.to_string()),
    err_5: (r#" "bar1" 'foo dar"#, r#"invalid string format (missing ' terminator)
|  "bar1" 'foo dar
| ...............^"#.to_string()),
    err_6: ("155 -44", r#"illegal char -
| 155 -44
| ....^"#.to_string()),
    err_7: ("255 1,3", r#"invalid number format (, is not a digit)
| 255 1,3
| .....^"#.to_string()),
    err_8: ("355 1e3", r#"invalid number format (e is not a digit)
| 355 1e3
| .....^"#.to_string()),
    err_9: ("455 0xEA", r#"invalid number format (x is not a digit)
| 455 0xEA
| .....^"#.to_string()),
    err_10: ("555 2,3-4", r#"invalid number format (, is not a digit)
| 555 2,3-4
| .....^"#.to_string()),
    err_11: ("655 3333_444_333", r#"invalid number format (_ is not a digit)
| 655 3333_444_333
| ........^"#.to_string()),
    err_12: ("755 33 44 -.3", r#"illegal char -
| 755 33 44 -.3
| ..........^"#.to_string()),
    err_13: ("foo:bar", r#"invalid REPLACE operation format (b is not equal char)
| foo:bar
| ....^"#.to_string()),
    err_14: ("foo:", r#"invalid REPLACE operation format (EOF)
| foo:
| ...^"#.to_string()),
    err_15: ("foo=", r#"invalid EQFILTER operation format (EOF)
| foo=
| ...^"#.to_string()),
}
