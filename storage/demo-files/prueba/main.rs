fn greet(name: &str) -> String {
    format!("Hola, {}", name)
}

fn main() {
    let languages = vec!["Rust", "Go", "PHP"];

    for language in languages {
        println!("{}", greet(language));
    }
}
