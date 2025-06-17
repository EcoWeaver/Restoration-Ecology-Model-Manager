<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		getBezierPath,
		MarkerType,
		Position,
		useSvelteFlow,
		type EdgeProps,
	} from "@xyflow/svelte";
	import type { Edge } from "$lib/model";
	import { relationElements } from "$lib/api";
	import { Trash2 } from "@lucide/svelte";

	let { id, data, sourceX, sourceY, targetX, targetY, markerEnd }: EdgeProps<Edge> = $props();

	let [path, labelX, labelY] = $derived(
		getBezierPath({
			sourceX,
			sourceY,
			sourcePosition: Position.Right,
			targetX,
			targetY,
			targetPosition: Position.Left,
		}),
	);

	const { updateEdge, getZoom, deleteElements } = useSvelteFlow();

	let value = $state(data!.kind);

	$effect(() => {
		const v = value;
		relationElements.then((re) => {
			updateEdge(
				id,
				{
					markerEnd: re.causal.includes(v) ? { type: MarkerType.ArrowClosed } : undefined,
					data: { kind: value, editable: data!.editable },
				},
				{ replace: true },
			);
		});
	});

	function deleteSelf() {
		deleteElements({ edges: [{ id }] });
	}
</script>

<BaseEdge {id} {path} {markerEnd} />
<EdgeLabel x={labelX} y={labelY}>
	{#if data!.editable}
		<div class="flex text-base">
			<select
				bind:value
				class="field-sizing-content h-10 cursor-pointer rounded-full bg-neutral-200 px-2 text-center dark:bg-neutral-700"
			>
				{#await relationElements then es}
					{#each Object.entries(es) as [key, values]}
						<optgroup label={key} style:font-size={`${getZoom()}em`}>
							{#each values as value}
								<option {value}>{value}</option>
							{/each}
						</optgroup>
					{/each}
				{/await}
			</select>
			<button
				class="ml-1 cursor-pointer rounded-full bg-neutral-200 p-2 dark:bg-neutral-700"
				onclick={deleteSelf}
				title="Delete"><Trash2 /></button
			>
		</div>
	{:else}
		<div
			class="flex h-10 items-center rounded-full bg-neutral-200 px-2 text-base dark:bg-neutral-700"
		>
			{data!.kind}
		</div>
	{/if}
</EdgeLabel>
