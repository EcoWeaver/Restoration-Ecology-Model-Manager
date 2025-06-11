<script lang="ts">
	let { onclick }: { onclick: () => any } = $props();
</script>

<button
	{onclick}
	class="m-4 ml-0 flex h-10 cursor-pointer place-items-center rounded-r-md border border-l-0 border-neutral-600 px-4 hover:bg-neutral-200"
>
	<span class="font-medium">search</span>
</button>
