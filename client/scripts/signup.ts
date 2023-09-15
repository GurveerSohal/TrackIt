import Cookies from "js-cookie";
import { server_url } from "../app";

const form = document.getElementById("signup-form");
form?.addEventListener("submit", signupUser);

async function signupUser(event) {
  event.preventDefault();
  const usernameInput = document.getElementById(
    "username-input"
  ) as HTMLInputElement;
  const passwordInput = document.getElementById(
    "password-input"
  ) as HTMLInputElement;

  const username = usernameInput.value;
  const password = passwordInput.value;
  const formData = {
    username: username,
    password: password,
  };

  const res = await fetch(`${server_url}/api/signup/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(formData),
  });

  console.log(res)

  if (!res.ok) {
    alert("Failed to login!")
    return
  }

  const data = await res.json()
  if (!data || !data.token) {
    alert("No token in response!")
    return
  }

  Cookies.set("token", data.token)
  console.log("cookie set")

  window.location.replace('/home')
}
