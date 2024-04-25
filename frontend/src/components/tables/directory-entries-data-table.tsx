"use client"
import { DirectoryEntries, DirectoryEntry, client } from "@/lib/contacts-api"
import { DataTable } from "./data-table"
import { PaginationState, createColumnHelper } from "@tanstack/react-table"
import { cookies, getCookieValue } from "@/lib/cookies"

import { OctagonAlert } from "lucide-react"

const ch = createColumnHelper<DirectoryEntry>();

const columns = [
	ch.accessor('entry_first_name', {
		header: 'First Name'
	}),
	ch.accessor('entry_last_name', {
		header: 'Last Name'
	}),
	ch.accessor('entry_email', {
		cell: (props) => {
			return (
				<span className="flex items-center gap-2">
					{props.row.original.entry_email_valid ? null : <OctagonAlert className="w-4 h-4" color="red" />}{props.getValue()}
				</span>
			)
		},
		header: 'Email'
	})
]

async function getData(id: number, props: PaginationState) {
	const res = await client.GET("/directories/{id}/entries", {
		params: {
			query: { limit: props.pageSize, offset: props.pageIndex },
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

	return res.data.entries
}

export const DirectoryEntriesDataTable = ({ id, data, rowCount }: { id: number, data: DirectoryEntries, rowCount: number }) => {
	return <DataTable queryKey={'entries'} defaultPageSize={100}
		pageSizes={[25, 50, 100, 250, 500]}
		data={data} columns={columns}
		rowCount={rowCount}
		queryFn={(state) => () => getData(id, state)} />
}
