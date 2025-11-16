const form = document.querySelector("#loginForm");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");
const adminCodeInput = document.querySelector("#adminCode");
const fileInput = document.querySelector("#file");

const messageElement = document.querySelector("#message");

form.addEventListener("submit", async (e) => {
	e.preventDefault();
	displayMessage("");

	const file = fileInput.files[0];
	const fileReader = new FileReader();
	fileReader.readAsDataURL(file);
	const dataUrl = await new Promise((resolve, reject) => {
		fileReader.onload = () => {
			resolve(fileReader.result);
		};
		fileReader.onerror = reject;
	});
	if (!dataUrl.includes(";base64")) {
		displayMessage(`Error, couldn't base64 encode file.`);
		return;
	}
	const content = dataUrl.slice(dataUrl.indexOf(",") + 1);
	const mime = dataUrl.slice("data:".length, dataUrl.indexOf(";"));

	const resp = await fetch("/api/v1/users/register-or-update", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			Authorization: `AdminCode ${adminCodeInput.value}`,
		},
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value,
			content,
			filename: file.name,
			mime,
		}),
	});
	if (!resp.ok) {
		displayMessage(
			`Error, received HTTP error code: ${
				resp.status
			}\nContent:\n${await resp.text()}`
		);
		return;
	}

	displayMessage("Success");
});

function displayMessage(message) {
	messageElement.innerText = message;
	messageElement.innerHTML = messageElement.innerHTML
		.split("\n")
		.join("<br>");

	messageElement.hidden = false;
}
