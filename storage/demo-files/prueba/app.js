const app = {
  name: "Axolotl",
  version: 1,
  enabled: true,
};

function greet(name) {
  return "Hola, " + name;
}

const languages = ["JavaScript", "Go", "PHP"];
languages.forEach((language) => {
  console.log(greet(language));

});