from dataclasses import dataclass


@dataclass
class User:
    name: str
    active: bool = True


def greet(user: User) -> str:
    status = "activo" if user.active else "inactivo"
    return f"Hola {user.name}, estado: {status}"


user = User("Axolotl")
print(greet(user))
