import type { Node as FlowNode, Edge as FlowEdge } from "@xyflow/svelte";
import type { Node as ApiNode, Relation as ApiRelation, Graph } from "./api";

export type Node = FlowNode<Omit<ApiNode, "position">, "custom">;
export type Edge = FlowEdge<Omit<ApiRelation, "from" | "to">, "custom">;

export const defaultEdgeOptions: Partial<Edge> = {
	type: "custom",
	data: {
		kind: "affects",
		position: { r: 0, phi: 0 },
	},
};

export function apiToFlow(
	graph: Graph,
	editable: boolean,
	idPrefix: string = "",
): { nodes: Node[]; edges: Edge[] } {
	const nodes: Node[] = Array.from(
		graph.nodes?.entries().map(([id, { body, position, color }]) => ({
			id: `${idPrefix}${id}`,
			type: "custom",
			dragHandle: ".drag-handle",
			position,
			origin: [0.5, 0.5],
			data: { body, color, editable },
		})) ?? [],
	);

	const edges: Edge[] =
		graph.relations?.map(({ from, to, kind, position }) => ({
			id: `${idPrefix}${from}-${idPrefix}${to}`,
			type: "custom",
			source: `${idPrefix}${from}`,
			target: `${idPrefix}${to}`,
			data: { kind, editable, position },
		})) ?? [];

	return { nodes, edges };
}

export function flowToApi(nodes: Node[], edges: Edge[]): Graph {
	const n: ApiNode[] = nodes.map((n) => {
		const { body, color } = n.data;
		return {
			body,
			color,
			position: n.position,
		};
	});

	const relations: ApiRelation[] = edges.map((e) => {
		const from = nodes.findIndex((n) => n.id === e.source);
		const to = nodes.findIndex((n) => n.id === e.target);

		const { kind, position } = e.data ?? defaultEdgeOptions.data!;
		if (position.r < 0) {
			position.r *= -1;
			position.phi += Math.PI;
		}

		// modulo (euclidean remainder)
		position.phi %= 2 * Math.PI;
		if (position.phi < 0) {
			position.phi += 2 * Math.PI;
		}

		return { from, to, kind, position };
	});

	return {
		nodes: n,
		relations,
	};
}
