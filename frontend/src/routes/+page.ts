import { maybeGoToSetup } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, fetch }) => {
	// TODO: fix development by adding CORS support
	await maybeGoToSetup(fetch);
};
