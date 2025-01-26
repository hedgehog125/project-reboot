const form = document.querySelector("#loginForm");
const usernameInput = document.querySelector("#username");
const passwordInput = document.querySelector("#password");
const authorizationCodeInput = document.querySelector("#authorizationCode");

const messageElement = document.querySelector("#message");

form.addEventListener("submit", async (e) => {
	e.preventDefault();
	displayMessage("");

	const resp = await fetch("/api/v1/users/download", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({
			username: usernameInput.value,
			password: passwordInput.value,
			authorizationCode: authorizationCodeInput.value,
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

	const {
		authorizationCode: newAuthorizationCode,
		authorizationCodeValidAt,
		content,
		filename,
		mime,
	} = await resp.json();

	if (content) {
		download(content, filename, mime);
		displayMessage("File downloaded");
		return;
	}

	const asDate = new Date(authorizationCodeValidAt);
	displayMessage(
		// TODO: include time
		`Success! The following authorisation code will be valid on ${asDate.toLocaleDateString()}.\n${newAuthorizationCode}`
	);
});

function displayMessage(message) {
	messageElement.innerText = message;
	messageElement.innerHTML = messageElement.innerHTML
		.split("\n")
		.join("<br>");

	messageElement.hidden = false;
}

function download(content, filename, mime) {
	// Ideally should use a Blob but this is good enough since it's in JSON anyway
	const url = `data:${mime};base64,${content}`;

	const anchor = document.createElement("a");
	anchor.href = url;
	anchor.download = filename;
	anchor.style.visibility = "none";
	document.body.appendChild(anchor);

	anchor.click();
	document.body.removeChild(anchor);
}
