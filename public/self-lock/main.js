const formPage1 = document.querySelector("#selfLockFormPage1");
const adminCodeInput = document.querySelector("#adminCode");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");
const untilInput = document.querySelector("#until");

const formPage2 = document.querySelector("#selfLockFormPage2");
const actionIDInput = document.querySelector("#actionID");
const twoFactorCodeInput = document.querySelector("#twoFactorCode");

const messageElement = document.querySelector("#message");

formPage1.addEventListener("submit", async (e) => {
	e.preventDefault();
	displayMessage("");

	const resp = await fetch("/api/v1/users/self-lock", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value,
			until: new Date(untilInput.value).toISOString(),
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

	const json = await resp.json();
	displayMessage(
		"Success. Enter the 2FA code below to confirm the self-lock."
	);
	formPage2.hidden = false;
	actionIDInput.value = json.twoFactorActionID;
});
formPage2.addEventListener("submit", async (e) => {
	e.preventDefault();
	displayMessage("");

	const resp = await fetch(
		`/api/v1/two-factor-actions/${actionIDInput.value}/confirm`,
		{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				code: twoFactorCodeInput.value,
			}),
		}
	);
	if (!resp.ok) {
		displayMessage(
			`Error, received HTTP error code: ${
				resp.status
			}\nContent:\n${await resp.text()}`
		);
		return;
	}

	displayMessage("Success.");
});

function displayMessage(message) {
	messageElement.innerText = message;
	messageElement.innerHTML = messageElement.innerHTML
		.split("\n")
		.join("<br>");

	messageElement.hidden = false;
}
