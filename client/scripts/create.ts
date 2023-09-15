import Cookies from "js-cookie";
import { server_url } from "../app";

let workout_number = -1;

const createWorkoutForm = document.getElementById("create-workout-form")
createWorkoutForm.addEventListener("submit", createSet)

const doneButton = document.getElementById("done-button");
doneButton.addEventListener("click", (event) => {
  event.preventDefault();
  window.location.href = '/home';
})

async function createWorkout() {
  const token = Cookies.get("token");
  if (!token) {
    window.location.replace("/login");
  }

  const res = await fetch(`${server_url}/api/create-workout/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ token }),
  });

  if (!res.ok) {
    alert("failed to create workout in database");
    console.log(res);
  }

  const data = await res.json();
  if (!data || !data.workout_number) {
    console.log("no data");
    console.log(res);
    return -1;
  }

  return data.workout_number;
}

async function createSet(event) {
  event.preventDefault();
  const token = Cookies.get("token");
  if (workout_number == -1) {
    workout_number = await createWorkout();
  }

  if (workout_number == -1) {
    alert("failed to create workout in database");
  }

  const name = (document.getElementById("name-input") as HTMLInputElement).value;
  const reps = parseInt((document.getElementById("reps-input") as HTMLInputElement).value);

  if (isNaN(reps)) {
    alert("enter an integer number of reps")
    return
  }

  const res = await fetch(`${server_url}/api/create-set/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ token, workout_number, name, reps}),
  });

  if (!res.ok) {
    alert('failed to create set')
    return
  }

  const data = await res.json()
  console.log(data)
}
