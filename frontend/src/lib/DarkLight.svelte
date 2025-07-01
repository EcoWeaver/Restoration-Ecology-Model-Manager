<script lang="ts">
	import { browser } from "$app/environment";
	import { Moon, Sun } from "@lucide/svelte";
	import { MediaQuery } from "svelte/reactivity";

	const preference = new MediaQuery("prefers-color-scheme: dark");
	let storage = $state(browser ? (localStorage?.getItem("theme") ?? "system") : null);
	let current = $derived(storage === "dark" || (storage === "system" && preference.current));

	function invertTheme(theme: string) {
		return theme === "dark" ? "light" : "dark";
	}

	function toggle() {
		const preferred = preference.current ? "dark" : "light";
		if (storage === "system") {
			storage = invertTheme(preferred);
		} else if (storage !== preferred) {
			storage = "system";
		} else {
			storage = invertTheme(storage);
		}
	}

	$effect(() => {
		if (storage) {
			localStorage.setItem("theme", storage);
		}

		document.documentElement.classList.toggle("dark", current);
	});
</script>

<button
	class="dark-light mr-1 cursor-pointer rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
	title="toggle dark/light theme"
	onclick={toggle}
>
	<span class="sun hidden"><Sun /></span>
	<span class="moon"><Moon /></span>
</button>

<style>
	:global(.dark) .dark-light .sun {
		display: unset;
	}

	:global(.dark) .dark-light .moon {
		display: none;
	}
</style>
