const form = document.querySelector("#loginForm");
const adminCodeInput = document.querySelector("#adminCode");
const usernameInput = document.querySelector("#username");
const discordInput = document.querySelector("#discordId");
const emailInput = document.querySelector("#email");

const messageElement = document.querySelector("#message");

form.addEventListener("submit", async (e) => {
	e.preventDefault();
	displayMessage("");

	const resp = await fetch("/api/v1/users/set-user-contacts", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			Authorization: `AdminCode ${adminCodeInput.value}`,
		},
		body: JSON.stringify({
			username: usernameInput.value,
			discordUserId: discordInput.value || null,
			email: emailInput.value || null,
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
		`Success. A test message was sent using these messengers: ${json.messagesSent.join(
			", "
		)}`
	);
});

function displayMessage(message) {
	messageElement.innerText = message;
	messageElement.innerHTML = messageElement.innerHTML
		.split("\n")
		.join("<br>");

	messageElement.hidden = false;
}
