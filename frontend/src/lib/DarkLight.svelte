<script lang="ts">
	import { browser } from "$app/environment";
	import { Moon, Sun } from "@lucide/svelte";
	import { MediaQuery } from "svelte/reactivity";

	let preference = new MediaQuery("prefers-color-scheme: dark");
	let storage = $state(browser ? localStorage?.getItem("theme") : null);
	let current = $derived(storage === "dark" || (storage === "system" && preference.current));

	function invertTheme(theme: string) {
		return theme === "dark" ? "light" : "dark";
	}

	function toggle() {
		let theme = storage ?? "system";
		let preferred = preference.current ? "dark" : "light";

		if (theme === "system") {
			storage = invertTheme(preferred);
		} else if (theme !== preferred) {
			storage = "system";
		} else {
			storage = invertTheme(theme);
		}
	}

	$effect(() => {
		if (storage) {
			localStorage.setItem("theme", storage);
		}

		document.documentElement.classList.toggle("dark", current);
	});
</script>

<svelte:head>
	<script>
		{
			const theme = localStorage.getItem("theme") ?? "system";
			document.documentElement.classList.toggle(
				"dark",
				theme === "dark" ||
					(theme === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches),
			);
		}
	</script>
</svelte:head>

<button class="mr-1 cursor-pointer" onclick={toggle}>
	{#if current}
		<Sun />
	{:else}
		<Moon />
	{/if}
</button>
