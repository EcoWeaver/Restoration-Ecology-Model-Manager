<script lang="ts">
	import { resolve } from "$app/paths";
	import { page } from "$app/state";
	import { getReasonPhrase } from "http-status-codes";

	const statusText = getReasonPhrase(page.status).toLowerCase();
</script>

<svelte:head>
	<title>{page.status} | {statusText}</title>
</svelte:head>

<div class="grid min-h-full items-center px-6 py-24 text-center sm:py-32 lg:px-8">
	<p class="text-2xl font-semibold tracking-tight text-red-500">
		{page.status}
		<span class="px-1">|</span>
		{statusText}
	</p>
	<h1 class="mt-4 text-5xl font-semibold tracking-tight text-balance wrap-anywhere sm:text-7xl">
		{page.error?.message || statusText}
	</h1>
	<div class="mt-10 items-center">
		<a
			class="rounded-md bg-red-300 px-5 py-2.5 text-lg font-semibold outline-offset-2 outline-red-400 hover:bg-red-400 focus-visible:outline-2 dark:bg-red-500"
			title="home"
			href={resolve("/")}
		>
			go home
		</a>
	</div>
</div>
