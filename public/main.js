const form = document.querySelector("#loginForm");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");
const authorizationCodeInput = document.querySelector("#authorizationCode");

const messageElement = document.querySelector("#message");

form.addEventListener("submit", async (e) => {
	e.preventDefault();

	const resp = await fetch("/v1/login", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value,
			authorizationCode: authorizationCodeInput.value,
		}),
	});
	if (!resp.ok) {
		displayMessage(`Error. Received HTTP error code: ${resp.status}`);
		return;
	}

	const {
		authorizationCode: newAuthorizationCode,
		authorizationCodeValidAt,
		rebootZipUrl,
	} = await resp.json();

	if (rebootZipUrl) {
		downloadUrl(rebootZipUrl);
		return;
	}

	const asDate = new Date(authorizationCodeValidAt);
	displayMessage(
		`Success! The following authorisation code will be valid on ${asDate.toLocaleDateString(
			undefined,
			{ timeStyle: "full" }
		)}.\n${newAuthorizationCode}`
	);
});

function displayMessage(message) {
	messageElement.innerText = message;
	messageElement.innerHTML = messageElement.innerHTML
		.split("\n")
		.join("<br>");

	messageElement.hidden = false;
}

function downloadUrl(url) {}
