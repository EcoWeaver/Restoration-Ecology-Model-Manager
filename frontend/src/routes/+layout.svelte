<script lang="ts">
	import "../app.css";
	import "../svelte-flow.css";

	let { children } = $props();
</script>

{@render children()}
