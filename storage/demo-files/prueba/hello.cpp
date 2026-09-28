#include <iostream>
#include <string>
#include <vector>

std::string greet(const std::string& name) {
    return "Hola, " + name;
}

int main() {
    std::vector<std::string> languages = {"C++", "Go", "PHP"};

    for (const auto& language : languages) {
        std::cout << greet(language) << '\n';
    }

    return 0;
}
