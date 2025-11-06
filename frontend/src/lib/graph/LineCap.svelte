<script lang="ts">
	const { id }: { id: string } = $props();
</script>

<svg>
	<marker {id} viewBox="-1.5 -1.5 3 3">
		<circle class="fill-(--xy-edge-stroke)" r="0.5" />
	</marker>
</svg>
