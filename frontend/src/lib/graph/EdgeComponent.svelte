<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		getStraightPath,
		MarkerType,
		Position,
		useInternalNode,
		useSvelteFlow,
		type EdgeProps,
		type InternalNode,
		type XYPosition,
	} from "@xyflow/svelte";
	import type { Edge } from "$lib/model";
	import { relationElements } from "$lib/api";
	import { Trash2 } from "@lucide/svelte";

	let { id, data, source, target, markerEnd }: EdgeProps<Edge> = $props();

	const sourceNode = useInternalNode(source);
	const targetNode = useInternalNode(target);

	function nodeIntersection(intersectionNode: InternalNode, targetNode: InternalNode) {
		const w = intersectionNode.measured.width! / 2;
		const h = intersectionNode.measured.height! / 2;

		const x2 = intersectionNode.internals.positionAbsolute.x + w;
		const y2 = intersectionNode.internals.positionAbsolute.y + h;
		const x1 = targetNode.internals.positionAbsolute.x + targetNode.measured.width! / 2;
		const y1 = targetNode.internals.positionAbsolute.y + targetNode.measured.height! / 2;

		const xx1 = (x1 - x2) / (2 * w) - (y1 - y2) / (2 * h);
		const yy1 = (x1 - x2) / (2 * w) + (y1 - y2) / (2 * h);
		const a = 1 / (Math.abs(xx1) + Math.abs(yy1));
		const xx3 = a * xx1;
		const yy3 = a * yy1;
		const x = w * (xx3 + yy3) + x2;
		const y = h * (-xx3 + yy3) + y2;

		return { x, y };
	}

	function edgePosition(node: InternalNode, intersectionPoint: XYPosition) {
		const n = { ...node.internals.positionAbsolute, ...node };
		const nx = Math.round(n.x);
		const ny = Math.round(n.y);
		const px = Math.round(intersectionPoint.x);
		const py = Math.round(intersectionPoint.y);

		if (px <= nx + 1) {
			return Position.Left;
		}
		if (px >= nx + n.measured.width! - 1) {
			return Position.Right;
		}
		if (py <= ny + 1) {
			return Position.Top;
		}
		if (py >= n.y + n.measured.height! - 1) {
			return Position.Bottom;
		}

		return Position.Top;
	}

	function edgeParams(source: InternalNode, target: InternalNode) {
		const sourceIntersectionPoint = nodeIntersection(source, target);
		const targetIntersectionPoint = nodeIntersection(target, source);

		const sourcePosition = edgePosition(source, sourceIntersectionPoint);
		const targetPosition = edgePosition(target, targetIntersectionPoint);

		return {
			sourceX: sourceIntersectionPoint.x,
			sourceY: sourceIntersectionPoint.y,
			targetX: targetIntersectionPoint.x,
			targetY: targetIntersectionPoint.y,
			sourcePosition,
			targetPosition,
		};
	}

	let [path, labelX, labelY] = $derived.by(() => {
		if (sourceNode.current && targetNode.current) {
			return getStraightPath(edgeParams(sourceNode.current, targetNode.current));
		}
		return [undefined!, undefined, undefined];
	});

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
<EdgeLabel x={labelX} y={labelY} class="!text-base">
	{#if data!.editable}
		<div class="flex">
			<select
				bind:value
				title={value}
				class="field-sizing-content h-10 cursor-pointer rounded-full bg-neutral-200 px-2 text-center outline-indigo-400 focus-visible:outline-2 dark:bg-neutral-700"
			>
				{#await relationElements then relationElements}
					{#each Object.entries(relationElements) as [group, values] (group)}
						<optgroup label={group} style:font-size={`${getZoom()}em`}>
							{#each values as value (value)}
								<option {value}>{value}</option>
							{/each}
						</optgroup>
					{/each}
				{/await}
			</select>
			<button
				class="ml-1 cursor-pointer rounded-full bg-neutral-200 p-2 outline-indigo-400 focus-visible:outline-2 dark:bg-neutral-700"
				onclick={deleteSelf}
				title="delete"><Trash2 /></button
			>
		</div>
	{:else}
		<div class="flex h-10 items-center rounded-full bg-neutral-200 px-2 dark:bg-neutral-700">
			{data!.kind}
		</div>
	{/if}
</EdgeLabel>
