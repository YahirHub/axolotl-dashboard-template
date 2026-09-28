import java.util.List;

public class Main {
    public static String greet(String name) {
        return "Hola, " + name;
    }

    public static void main(String[] args) {
        List<String> languages = List.of("Java", "Go", "PHP");

        for (String language : languages) {
            System.out.println(greet(language));
        }
    }
}
