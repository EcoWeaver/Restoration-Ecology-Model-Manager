import { apiBase } from "#lib/api.js";
import { redirect } from "@sveltejs/kit";
import { StatusCodes } from "http-status-codes";

export function load() {
	redirect(StatusCodes.PERMANENT_REDIRECT, `${apiBase}/swagger/index.html`, {
		external: [apiBase],
	});
}
