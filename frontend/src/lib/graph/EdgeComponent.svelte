<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		MarkerType,
		useInternalNode,
		useStore,
		useSvelteFlow,
		type EdgeProps,
		type InternalNode,
	} from "@xyflow/svelte";
	import type { Edge } from "$lib/model";
	import { cartesianToPolar, useElements, polarToCartesian, type Cartesian } from "$lib/api";
	import { Move, RotateCcw, Trash2 } from "@lucide/svelte";
	import LineCap from "./LineCap.svelte";
	import { untrack } from "svelte";
	import { slide } from "svelte/transition";

	let { id, data, source, target, markerStart, markerEnd }: EdgeProps<Edge> = $props();

	const { getRelationElements } = useElements();
	const { updateEdge, getZoom, deleteElements, screenToFlowPosition } = useSvelteFlow();
	const sourceNode = $derived(useInternalNode(source));
	const targetNode = $derived(useInternalNode(target));
	const { edgesFocusable } = $derived(useStore());

	function nodeIntersection(intersectionNode: InternalNode, target: Cartesian): Cartesian {
		// https://math.stackexchange.com/questions/1724792/an-algorithm-for-finding-the-intersection-point-between-a-center-of-vision-and-a

		const w = intersectionNode.measured.width! / 2;
		const h = intersectionNode.measured.height! / 2;
		const x2 = intersectionNode.internals.positionAbsolute.x + w;
		const y2 = intersectionNode.internals.positionAbsolute.y + h;
		const x1 = target.x;
		const y1 = target.y;

		const x1_ = (x1 - x2) / (2 * w) - (y1 - y2) / (2 * h);
		const y1_ = (x1 - x2) / (2 * w) + (y1 - y2) / (2 * h);
		const a = 1 / (Math.abs(x1_) + Math.abs(y1_));
		const x3 = a * x1_;
		const y3 = a * y1_;
		const x = w * (x3 + y3) + x2;
		const y = h * (-x3 + y3) + y2;

		const cr = 8;
		const offx = w - cr;
		const offy = h - cr;

		if ((x >= x2 - offx && x <= x2 + offx) || (y >= y2 - offy && y <= y2 + offy)) {
			return { x, y };
		}

		const ccx = x < x2 - offx ? x2 - offx : x2 + offx;
		const ccy = y < y2 - offy ? y2 - offy : y2 + offy;

		// https://www.bluebill.net/circle_ray_intersection.html

		let dx = x1 - x2;
		let dy = y1 - y2;
		const dl = Math.hypot(dx, dy);
		dx /= dl;
		dy /= dl;

		const ux = ccx - x2;
		const uy = ccy - y2;

		const u1d = ux * dx + uy * dy;
		const u1x = u1d * dx;
		const u1y = u1d * dy;

		const u2x = ux - u1x;
		const u2y = uy - u1y;

		const d2 = u2x * u2x + u2y * u2y;
		const m = Math.sqrt(cr * cr - d2);

		return {
			x: x2 + u1x + m * dx,
			y: y2 + u1y + m * dy,
		};
	}

	function controlPoints(start: number, middle: number, end: number): number[] {
		const x = (end - start) / 6;

		const first1 = (start + middle) / 2 - x / 2;
		const first2 = middle - x;

		const second1 = middle + x;
		const second2 = (end + middle) / 2 + x / 2;

		return [first1, first2, second1, second2];
	}

	let [path, labelX, labelY, middleX, middleY, alpha] = $derived.by(() => {
		const source = sourceNode.current!;
		const target = targetNode.current!;

		const sourceX = source.internals.positionAbsolute.x + source.measured.width! / 2;
		const sourceY = source.internals.positionAbsolute.y + source.measured.height! / 2;
		const targetX = target.internals.positionAbsolute.x + target.measured.width! / 2;
		const targetY = target.internals.positionAbsolute.y + target.measured.height! / 2;

		const middleX = (sourceX + targetX) / 2;
		const middleY = (sourceY + targetY) / 2;

		const alpha = Math.atan2(targetY - sourceY, targetX - sourceX);

		const offset = polarToCartesian({
			r: data!.position.r,
			phi: data!.position.phi + alpha,
		});

		const labelX = middleX + offset.x;
		const labelY = middleY + offset.y;

		const { x: edgeSourceX, y: edgeSourceY } = nodeIntersection(source, { x: labelX, y: labelY });
		const { x: edgeTargetX, y: edgeTargetY } = nodeIntersection(target, { x: labelX, y: labelY });

		const x = controlPoints(edgeSourceX, labelX, edgeTargetX);
		const y = controlPoints(edgeSourceY, labelY, edgeTargetY);

		return [
			`M ${edgeSourceX},${edgeSourceY} C ${x[0]},${y[0]} ${x[1]},${y[1]} ${labelX},${labelY} C ${x[2]},${y[2]} ${x[3]},${y[3]} ${edgeTargetX},${edgeTargetY}`,
			labelX,
			labelY,
			middleX,
			middleY,
			alpha,
		];
	});

	let { kind, position } = $derived(data!);
	$effect(() => {
		const kind_ = kind;
		const position_ = position;

		const edge = untrack(() => ({
			markerStart: `linecap-s-${source}-${target}`,
			markerEnd: getRelationElements().causal.includes(kind)
				? { type: MarkerType.ArrowClosed, color: null }
				: `linecap-e-${source}-${target}`,
			data: { kind: kind_, position: position_ },
		}));
		updateEdge(id, edge, { replace: true });
	});

	const drag: {
		dragging: boolean;
		offset?: Cartesian;
	} = {
		dragging: false,
		offset: undefined,
	};

	function onpointermove(e: PointerEvent) {
		const pos = screenToFlowPosition({ x: e.pageX, y: e.pageY });
		if (!drag.offset) {
			drag.offset = { x: pos.x - labelX, y: pos.y - labelY };
		}
		position = cartesianToPolar({
			x: pos.x - drag.offset.x - middleX,
			y: pos.y - drag.offset.y - middleY,
		});
		position.phi -= alpha;
	}

	function deleteSelf() {
		deleteElements({ edges: [{ id }] });
	}
</script>

<LineCap id="linecap-s-{source}-{target}" />

{#if !getRelationElements().causal.includes(kind)}
	<LineCap id="linecap-e-{source}-{target}" />
{/if}

<BaseEdge {id} {path} {markerStart} {markerEnd} />
<EdgeLabel x={labelX} y={labelY} class="text-base!">
	<div class="flex">
		{#if edgesFocusable}
			<select
				bind:value={kind}
				title={kind}
				class="field-sizing-content h-10 cursor-pointer rounded-full bg-neutral-200 px-2 text-center outline-indigo-400 focus-visible:outline-2 dark:bg-neutral-700"
			>
				{#each Object.entries(getRelationElements()) as [group, values] (group)}
					<optgroup label={group} style:font-size="{getZoom()}em">
						{#each values as value (value)}
							<option {value}>{value}</option>
						{/each}
					</optgroup>
				{/each}
			</select>
		{:else}
			<div class="flex h-10 items-center rounded-full bg-neutral-200 px-2 dark:bg-neutral-700">
				{data!.kind}
			</div>
		{/if}

		{#if edgesFocusable}
			<div
				class="nopan ml-1 flex h-10 cursor-default gap-2 rounded-full bg-neutral-200 p-2 dark:bg-neutral-700"
				transition:slide={{ duration: 100, axis: "x" }}
			>
				<div
					class="inline-block cursor-grab"
					onpointerdown={() => {
						drag.dragging = true;
						drag.offset = undefined;
					}}
				>
					<Move />
				</div>
				<button
					class="cursor-pointer rounded-full outline-offset-3 outline-indigo-400 focus-visible:outline-2"
					onclick={() => (position = { r: 0, phi: 0 })}
					title="reset position"
				>
					<RotateCcw />
				</button>
				<button
					class="cursor-pointer rounded-full outline-offset-3 outline-indigo-400 focus-visible:outline-2"
					onclick={deleteSelf}
					title="delete"
				>
					<Trash2 />
				</button>
			</div>
		{/if}
	</div>
</EdgeLabel>

<svelte:window
	onpointermove={drag.dragging ? onpointermove : undefined}
	onpointerup={drag.dragging ? () => (drag.dragging = false) : undefined}
/>
