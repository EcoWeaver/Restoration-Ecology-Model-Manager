import { env } from "$env/dynamic/public";
import { redirect } from "@sveltejs/kit";
import { StatusCodes } from "http-status-codes";

export function load() {
	const apiUrl = env.PUBLIC_API_URL || "http://localhost:8080";
	redirect(StatusCodes.PERMANENT_REDIRECT, `${apiUrl}/swagger/index.html`);
}
