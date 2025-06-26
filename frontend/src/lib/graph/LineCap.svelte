<script lang="ts">
	const { id }: { id: string } = $props();
</script>

<svg>
	<marker {id} viewBox="0 0 6 6" refX={1.5} refY={1.5}>
		<circle class="svelte-flow__remma-custom-linecap" cx="1.5" cy="1.5" r="0.5" />
	</marker>
</svg>
