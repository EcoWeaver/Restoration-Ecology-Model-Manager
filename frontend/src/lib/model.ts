import type { Node as FlowNode, Edge as FlowEdge } from "@xyflow/svelte";
import {
	putModel,
	type Node as ApiNode,
	type Relation as ApiRelation,
	type Graph,
	type Metadata,
	type SaveResult,
	type UnsavedModel,
} from "./api";

export type Node = FlowNode<Omit<ApiNode, "position"> & { editable: boolean }, "custom"> & {
	dragHandle: ".drag-handle";
};
export type Edge = FlowEdge<Omit<ApiRelation, "from" | "to"> & { editable: boolean }, "custom">;

export function defaultEdge(editable: boolean): Partial<Edge> {
	return {
		type: "custom",
		data: {
			kind: "correlated-with",
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

export async function saveModel(
	parent: string | undefined,
	metadata: Metadata[],
	n: Node[],
	e: Edge[],
): Promise<SaveResult> {
	const nodes: ApiNode[] = n.map((n) => {
		const { body, color } = n.data;
		return {
			body,
			color,
			position: n.position,
		};
	});

	const relations: ApiRelation[] = e.map((e) => {
		const from = n.findIndex((n) => n.id === e.source);
		const to = n.findIndex((n) => n.id === e.target)!;

		const { kind, position } = e.data!;
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

	const graph = { nodes, relations };
	const request: UnsavedModel = {
		metadata,
		graph,
		parent,
	};

	return await putModel(request);
}
