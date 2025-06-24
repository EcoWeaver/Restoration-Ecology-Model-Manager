import type { Node as FlowNode, Edge as FlowEdge } from "@xyflow/svelte";
import {
	putModel,
	type Node as ApiNode,
	type Relation as ApiRelation,
	type Graph,
	type Metadata,
	type ModelRequest,
} from "./api";

export type Node = FlowNode<Omit<ApiNode, "position"> & { editable: boolean }, "custom"> & {
	dragHandle: ".drag-handle";
};
export type Edge = FlowEdge<Omit<ApiRelation, "from" | "to"> & { editable: boolean }, "custom">;

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
		graph.relations?.map(({ from, to, kind }) => ({
			source: from.toString(),
			target: to.toString(),
			type: "custom",
			id: `${from}-${to}`,
			data: { kind, editable },
		})) ?? [];

	return { nodes, edges };
}

export async function saveModel(
	identifier: string,
	metadata: Metadata[],
	n: Node[],
	e: Edge[],
	version?: string,
): Promise<ModelRequest> {
	const nodes: ApiNode[] = n.map((n) => {
		const { body, color } = n.data;
		return {
			body,
			color,
			position: {
				x: Math.round(n.position.x),
				y: Math.round(n.position.y),
			},
		};
	});

	const relations: ApiRelation[] = e.map((e) => {
		const from = n.findIndex((n) => n.id === e.source);
		const to = n.findIndex((n) => n.id === e.target)!;

		const { kind } = e.data!;
		return { kind, from, to };
	});

	const graph = { nodes, relations };
	const request: ModelRequest = {
		identifier,
		version: version ?? "",
		metadata,
		graph,
	};

	return await putModel(request);
}
