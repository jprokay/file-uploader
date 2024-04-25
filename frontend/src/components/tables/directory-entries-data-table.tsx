"use client"
import { DirectoryEntries, DirectoryEntry, client } from "@/lib/contacts-api"
import { DataTable } from "./data-table"
import { ColumnDef, PaginationState, createColumnHelper } from "@tanstack/react-table"
import { cookies, getCookieValue } from "@/lib/cookies"

import { OctagonAlert } from "lucide-react"

const ch = createColumnHelper<DirectoryEntry>();

const columns: ColumnDef<DirectoryEntry>[] = [
	{
		accessorKey: 'order_id',
		header: 'Row'
	},
	{
		accessorKey: 'entry_first_name',
		header: 'First Name'
	},

	{
		accessorKey: 'entry_last_name',
		header: 'Last Name'
	},
	{
		accessorKey: 'entry_email',
		header: 'Email',
		cell: (props) => {
			return (
				<span className="flex items-center gap-2">
					{props.row.original.entry_email_valid ? null : <OctagonAlert className="w-4 h-4" color="red" />}{String(props.getValue())}
				</span>
			)
		},
	}
]

async function getData(id: number, props: PaginationState, search: string | undefined) {
	const res = await client.GET("/directories/{id}/entries", {
		params: {
			query: { limit: props.pageSize, offset: props.pageIndex, search },
			cookie: { userId: getCookieValue("userId") || "" },
			path: {
				id
			}
		},
		headers: { Cookie: cookies() },
	})

	if (res.error) {
		throw new Error('Oops something went wrong')
	}

	return res.data
}

export const DirectoryEntriesDataTable = ({ id }: { id: number }) => {
	return <DataTable queryKey={'entries'} defaultPageSize={100}
		pageSizes={[25, 50, 100, 250, 500]}
		columns={columns}
		enableSearch={true}
		queryFn={(state, search) => () => getData(id, state, search)} />
}
