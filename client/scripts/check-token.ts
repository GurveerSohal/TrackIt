import Cookies from "js-cookie";
import { server_url } from "../app";

async function checkToken() {
    const token = Cookies.get("token")
    if (!token) {
        window.location.replace('/login')
    }


  const res = await fetch(`${server_url}/api/token/verify/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({token}),
  });

  if (!res.ok) {
    window.location.replace('/login')
    return
  }

  console.log("token is valid")
}

checkToken()