import type { Node as FlowNode, Edge as FlowEdge } from "@xyflow/svelte";
import { type Node as ApiNode, type Relation as ApiRelation, type Graph } from "./api";

export type Node = FlowNode<Omit<ApiNode, "position"> & { editable: boolean }, "custom"> & {
	dragHandle: ".drag-handle";
};
export type Edge = FlowEdge<Omit<ApiRelation, "from" | "to"> & { editable: boolean }, "custom">;

export function defaultEdge(editable: boolean): Partial<Edge> {
	return {
		type: "custom",
		data: {
			kind: "affects",
			editable,
			position: { r: 0, phi: 0 },
		},
	};
}

export function apiToFlow(graph: Graph, editable: boolean): { nodes: Node[]; edges: Edge[] } {
	const nodes: Node[] = Array.from(
		graph.nodes?.entries().map(([id, { body, position, color }]) => ({
			position,
			type: "custom",
			dragHandle: ".drag-handle",
			id: id.toString(),
			data: { body, color, editable },
		})) ?? [],
	);

	const edges: Edge[] =
		graph.relations?.map(({ from, to, kind, position }) => ({
			source: from.toString(),
			target: to.toString(),
			type: "custom",
			id: `${from}-${to}`,
			data: { kind, editable, position },
		})) ?? [];

	return { nodes, edges };
}

export function toModel(nodes: Node[], edges: Edge[]): Graph {
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

		const { kind, position } = e.data ?? defaultEdge(true).data!;
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
