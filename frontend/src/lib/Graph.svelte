<script lang="ts">
	import type { Snippet } from "svelte";
	import {
		Background,
		BackgroundVariant,
		ConnectionLineType,
		ConnectionMode,
		Panel,
		SvelteFlow,
		useSvelteFlow,
		type OnConnectEnd,
		type OnMove,
		type XYPosition,
	} from "@xyflow/svelte";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import { type Node, type Edge, defaultEdge } from "$lib/model";
	import "@xyflow/svelte/dist/style.css";

	let {
		nodes = $bindable(),
		edges = $bindable(),
		addNode,
		onmoveend,
		children,
	}: {
		nodes: Node[];
		edges: Edge[];
		addNode?: (pos: XYPosition) => string;
		onmoveend?: OnMove;
		children: Snippet;
	} = $props();

	const { screenToFlowPosition } = useSvelteFlow();
	const editable = $derived(addNode !== undefined);

	const onconnectend: OnConnectEnd = (event, connectionState) => {
		if (!addNode || connectionState.isValid) {
			return;
		}

		const sourceId = connectionState.fromNode!.id;
		const sourceNode = nodes.find((n) => n.id === sourceId)!;
		const sourcePos = sourceNode.position;
		const sourceWidth = sourceNode.measured!.width! / 2 + 40;
		const sourceHeight = sourceNode.measured!.height! / 2 + 40;

		const { clientX, clientY } = "changedTouches" in event ? event.changedTouches[0] : event;
		const targetPos = screenToFlowPosition({ x: clientX, y: clientY });

		if (
			targetPos.x > sourcePos.x - sourceWidth &&
			targetPos.x < sourcePos.x + sourceWidth &&
			targetPos.y > sourcePos.y - sourceHeight &&
			targetPos.y < sourcePos.y + sourceHeight
		) {
			return;
		}

		const targetId = addNode(targetPos);

		const newEdge = {
			...defaultEdge,
			id: `${sourceId}-${targetId}`,
			source: sourceId,
			target: targetId,
		};
		edges = [...edges, newEdge];
	};
</script>

<SvelteFlow
	bind:nodes
	bind:edges
	nodeTypes={{ custom: NodeComponent }}
	edgeTypes={{ custom: EdgeComponent }}
	deleteKey={["Delete", "Backspace"]}
	fitView={nodes.length > 0}
	connectionMode={ConnectionMode.Loose}
	nodesDraggable={editable}
	nodesConnectable={editable}
	elementsSelectable={editable}
	nodesFocusable={editable}
	edgesFocusable={editable}
	defaultEdgeOptions={defaultEdge(editable)}
	connectionLineType={ConnectionLineType.Straight}
	elevateEdgesOnSelect={false}
	proOptions={{ hideAttribution: true }}
	{onmoveend}
	onconnectend={editable ? onconnectend : undefined}
>
	<Panel position="bottom-right">
		{@render children()}
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
