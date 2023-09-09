import Cookies from "js-cookie";
import { server_url } from "../app";

const workoutHistoryDiv = document.getElementById("workout-history");
async function getWorkouts() {
  const token = Cookies.get("token");
  if (!token) {
    window.location.replace("/login");
  }

  const res = await fetch(`${server_url}/api/workouts/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ token }),
  });

  if (!res.ok) {
    console.log("error when getting workouts");
    console.log(res);
    return;
  }

  const data = await res.json()
  if (!data || !data.workouts) {
    console.log('no data')
    console.log(res)
    return
  }

  workoutHistoryDiv.innerHTML = JSON.stringify(data)
}

getWorkouts()
