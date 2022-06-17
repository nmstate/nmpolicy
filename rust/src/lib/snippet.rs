pub(crate) fn snippet(expression: String, mut pos: usize) -> String {
    if expression.is_empty() {
        return "".to_string();
    }

    if pos >= expression.len() {
        pos = expression.len() - 1
    }

    let mut marker = String::from("");
    for _ in 0..pos {
        marker.push('.');
    }
    marker.push('^');
    return format!("| {}\n| {}", expression, marker);
}
