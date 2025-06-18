<script lang="ts">
	import { page } from "$app/state";
	import { getReasonPhrase } from "http-status-codes";

	let statusText = getReasonPhrase(page.status).toLowerCase();
</script>

<svelte:head>
	<title>{page.status} | {statusText}</title>
</svelte:head>

<main class="grid min-h-full place-items-center px-6 py-24 sm:py-32 lg:px-8">
	<div class="text-center">
		<p class="text-2xl font-semibold tracking-tight text-red-600">
			{page.status}
			<span class="px-1">|</span>
			{statusText}
		</p>
		<h1 class="mt-4 text-5xl font-semibold tracking-tight text-balance text-gray-900 sm:text-7xl">
			{page.error?.message || statusText}
		</h1>
		<div class="mt-10 items-center">
			<a
				href="/"
				class="rounded-md bg-red-600 px-5 py-2.5 text-lg font-semibold text-white hover:bg-red-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
				>go home</a
			>
		</div>
	</div>
</main>
