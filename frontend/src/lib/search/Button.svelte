<script lang="ts">
	let { onclick }: { onclick: () => any } = $props();
</script>

<button
	{onclick}
	class="m-4 ml-0 flex h-10 cursor-pointer place-items-center rounded bg-indigo-300 px-4 hover:bg-indigo-400"
>
	<span class="font-medium">search</span>
</button>
