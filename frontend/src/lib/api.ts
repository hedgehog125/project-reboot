import { PUBLIC_API_DOMAIN } from "$env/static/public";
import { page } from "$app/state";
import { goto } from "$app/navigation";

class StatusError extends Error {
	jsonResponse: JsonResponse;

	constructor(resp: JsonResponse) {
		super(`request failed with status ${resp.status}`);
		this.jsonResponse = resp;
	}
}

export class JsonResponse {
	headers: Headers;
	status: number;
	data: any;
	redirecting: boolean;

	constructor(resp: Response, data: any) {
		this.headers = resp.headers;
		this.status = resp.status;
		this.data = data;
		this.redirecting = false;
	}

	get ok(): boolean {
		return this.status >= 200 && this.status <= 299;
	}
	throwForStatus() {
		if (!this.ok) {
			throw new StatusError(this);
		}
	}
}

export async function fetchJson(
	fetch: typeof global.fetch,
	url: string,
	init?: RequestInit | undefined,
): Promise<JsonResponse> {
	const urlObj = new URL(PUBLIC_API_DOMAIN + url);
	const resp = await fetch(urlObj, init);
	const json = await resp.json();

	const jsonResponse = new JsonResponse(resp, json);
	if (
		resp.status === 404 &&
		responseHasErrorCode(jsonResponse, "ENDPOINT_NOT_FOUND") &&
		!page.route.id?.startsWith("/setup") &&
		!urlObj.pathname.startsWith("/api/v1/setup/")
	) {
		if (await maybeGoToSetup(fetch)) {
			jsonResponse.redirecting = true;
		}
	}

	return jsonResponse;
}

export function responseHasErrorCode(response: JsonResponse, errorCode: string): boolean {
	const errors = response.data?.errors;
	if (!Array.isArray(errors)) return false;

	return errors.find((error) => error?.code === errorCode) != null;
}

export async function maybeGoToSetup(fetch: typeof global.fetch): Promise<boolean> {
	const setupStatus = await fetchJson(fetch, "/api/v1/setup/");
	setupStatus.throwForStatus();
	if (setupStatus.data.isComplete) {
		return false;
	}
	// TODO: base URL support?
	if (!setupStatus.data.isEnvComplete) {
		goto("/setup/env");
	} else {
		goto("/setup/admin-messengers");
	}
	return true;
}
