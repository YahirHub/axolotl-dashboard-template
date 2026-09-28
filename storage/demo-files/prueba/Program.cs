using System;
using System.Collections.Generic;

public class Program
{
    static string Greet(string name) => $"Hola, {name}";

    public static void Main()
    {
        var languages = new List<string> { "C#", "Go", "PHP" };

        foreach (var language in languages)
        {
            Console.WriteLine(Greet(language));
        }
    }
}
