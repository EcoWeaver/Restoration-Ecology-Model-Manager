<script lang="ts">
	import {
		BaseEdge,
		EdgeLabel,
		getEdgeCenter,
		MarkerType,
		Position,
		useInternalNode,
		useSvelteFlow,
		type EdgeProps,
		type InternalNode,
		type XYPosition,
	} from "@xyflow/svelte";
	import type { Edge } from "$lib/model";
	import { cartesianToPolar, relationElements as getRelationElements } from "$lib/api";
	import { ArrowDown, ArrowUp, Move, Trash2 } from "@lucide/svelte";
	import LineCap from "./LineCap.svelte";

	let { id, data, source, target, markerStart, markerEnd }: EdgeProps<Edge> = $props();
	const relationElements = getRelationElements();

	const sourceNode = useInternalNode(source);
	const targetNode = useInternalNode(target);

	function nodeIntersection(intersectionNode: InternalNode, targetNode: InternalNode) {
		/* eslint-disable camelcase -- maths isn't camelcase */

		// https://math.stackexchange.com/questions/1724792/an-algorithm-for-finding-the-intersection-point-between-a-center-of-vision-and-a

		const w = intersectionNode.measured.width! / 2;
		const h = intersectionNode.measured.height! / 2;
		const x2 = intersectionNode.internals.positionAbsolute.x + w;
		const y2 = intersectionNode.internals.positionAbsolute.y + h;
		const x1 = targetNode.internals.positionAbsolute.x + targetNode.measured.width! / 2;
		const y1 = targetNode.internals.positionAbsolute.y + targetNode.measured.height! / 2;

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
		const dl = Math.sqrt(dx * dx + dy * dy);
		dx /= dl;
		dy /= dl;

		const ux = ccx - x2;
		const uy = ccy - y2;

		const u1d = ux * dx + uy * dy;
		const u1x = u1d * dx;
		const u1y = u1d * dy;

		const u2x = ux - u1x;
		const u2y = uy - u1y;

		const d_2 = u2x * u2x + u2y * u2y;
		const m = Math.sqrt(cr * cr - d_2);

		return {
			x: x2 + u1x + m * dx,
			y: y2 + u1y + m * dy,
		};
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

	function controlPoints(start: number, middle: number, end: number) {
		const x = (end - start) / 6;

		const first1 = (start + middle) / 2 - x / 2;
		const first2 = middle - x;

		const second1 = middle + x;
		const second2 = (end + middle) / 2 + x / 2;

		return [first1, first2, second1, second2];
	}

	let [path, labelX, labelY, middleX, middleY] = $derived.by(() => {
		if (sourceNode.current && targetNode.current) {
			const params = edgeParams(sourceNode.current, targetNode.current);

			const dX = params.targetX - params.sourceX;
			const dY = params.targetY - params.sourceY;

			const length = Math.hypot(dX, dY);

			const [middleX, middleY] = getEdgeCenter(params);

			const labelX = middleX + (dY / length) * data!.position!.r;
			const labelY = middleY - (dX / length) * data!.position!.r;

			const x = controlPoints(params.sourceX, labelX, params.targetX);
			const y = controlPoints(params.sourceY, labelY, params.targetY);

			return [
				`M ${params.sourceX},${params.sourceY} C ${x[0]},${y[0]} ${x[1]},${y[1]} ${labelX},${labelY} C ${x[2]},${y[2]} ${x[3]},${y[3]} ${params.targetX},${params.targetY}`,
				labelX,
				labelY,
				middleX,
				middleY,
			];
		}
		return [undefined!, undefined, undefined];
	});

	const { updateEdge, getZoom, deleteElements } = useSvelteFlow();

	let kind = $state(data!.kind);
	let position = $state(data!.position);

	$effect(() => {
		const kindState = kind;
		const positionState = position;
		relationElements.then((re) => {
			updateEdge(
				id,
				{
					markerStart: `linecap-s-${source}-${target}`,
					markerEnd: re.causal.includes(kindState)
						? { type: MarkerType.ArrowClosed }
						: `linecap-e-${source}-${target}`,
					data: { kind: kindState, editable: data!.editable, position: positionState },
				},
				{ replace: true },
			);
		});
	});

	const { screenToFlowPosition } = useSvelteFlow();
	const drag = $state({ dragging: false, offsetInvalid: true, offsetX: 0, offsetY: 0 });

	function deleteSelf() {
		deleteElements({ edges: [{ id }] });
	}
</script>

<LineCap id={`linecap-s-${source}-${target}`} />

{#await relationElements then re}
	{#if !re.causal.includes(kind)}
		<LineCap id={`linecap-e-${source}-${target}`} />
	{/if}
{/await}

<BaseEdge {id} {path} {markerStart} {markerEnd} />
<EdgeLabel x={labelX} y={labelY} class="!text-base">
	{#if data!.editable}
		<div class="flex">
			<select
				bind:value={kind}
				title={kind}
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
			<div class="ml-1 h-10 rounded-full bg-neutral-200 p-2 dark:bg-neutral-700">
				<button
					class="cursor-pointer rounded-4xl outline-offset-3 outline-indigo-400 focus-visible:outline-2"
					onclick={() => (position!.r += 25)}
					title="moveUp"
				>
					<ArrowUp />
				</button>
				<button
					class="cursor-pointer rounded-4xl outline-offset-3 outline-indigo-400 focus-visible:outline-2"
					onclick={() => (position!.r -= 25)}
					title="moveDown"
				>
					<ArrowDown />
				</button>
				<div
					class="nopan nodrag inline-block"
					class:cursor-grab={!drag.dragging}
					class:cursor-grabbing={drag.dragging}
					onpointerdown={() => {
						drag.dragging = true;
						drag.offsetInvalid = true;
					}}
					onpointerup={() => {
						drag.dragging = false;
					}}
					onpointermove={(e) => {
						if (drag.dragging) {
							const pos = screenToFlowPosition({ x: e.pageX, y: e.pageY });
							if (drag.offsetInvalid) {
								drag.offsetInvalid = false;
								drag.offsetX = pos.x - labelX!;
								drag.offsetY = pos.y - labelY!;
							}
							position = cartesianToPolar({
								x: pos.x - drag.offsetX - middleX!,
								y: pos.y - drag.offsetY - middleY!,
							});
						}
					}}
				>
					<Move />
				</div>
				<button
					class="cursor-pointer rounded-4xl outline-offset-3 outline-indigo-400 focus-visible:outline-2"
					onclick={deleteSelf}
					title="delete"
				>
					<Trash2 />
				</button>
			</div>
		</div>
	{:else}
		<div class="flex h-10 items-center rounded-full bg-neutral-200 px-2 dark:bg-neutral-700">
			{data!.kind}
		</div>
	{/if}
</EdgeLabel>
