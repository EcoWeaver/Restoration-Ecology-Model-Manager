<script lang="ts">
	import { X } from "@lucide/svelte";

	let open = $state(false);
	let text = $state("");
	let isError = $state(false);
	let timeoutId = $state<number | number>();

	export function error(err: string) {
		clearTimeout(timeoutId);

		text = err;
		open = true;
		isError = true;
		timeoutId = setTimeout(() => (open = false), 6000);
	}

	export function info(info: string) {
		clearTimeout(timeoutId);

		text = info;
		open = true;
		isError = false;
		timeoutId = setTimeout(() => (open = false), 3000);
	}
</script>

{#if open}
	<div role="alert" class="t-0 absolute z-1000 mt-18 flex w-full justify-center">
		<span
			class:bg-red-100={isError}
			class:text-red-700={isError}
			class:dark:text-red-100={isError}
			class:dark:bg-red-500={isError}
			class="flex items-center rounded-md bg-indigo-300 px-5 py-3 dark:bg-indigo-500"
		>
			<span>{text || "popup yay"}</span>
			<button class="cursor-pointer pl-2" onclick={() => (open = false)}><X /></button>
		</span>
	</div>
{/if}
