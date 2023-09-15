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

  const workouts = data.workouts
  console.log(workouts)

  let html = ""
  html += `<h1 class="text-primary page-heading">Total workouts: ${Object.keys(workouts).length}</h1>`

  for (const workout_number in workouts) {
    let workout_html = `
    <table>
      <tr>
        <th colspan="3">
          Workout Number: ${workout_number}
        </th>
      </tr>
    `
    workout_html += `
    <tr>
      <th>Set</th>
      <th>Name</th>
      <th>Reps</th>
    </tr>
    `
    const set_array = workouts[workout_number]
    set_array.forEach(set => {
      workout_html += `
      <tr>
        <td>${set.set_number}</td>
        <td>${set.name}</td>
        <td>${set.reps}</td>
      </tr>
      `
    });

    workout_html += `</table>`

    html += workout_html
  }

  workoutHistoryDiv.innerHTML = html
}

getWorkouts()
