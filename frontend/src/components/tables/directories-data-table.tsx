"use client"
import { Directories, Directory, client } from "@/lib/contacts-api"
import { DataTable } from "./data-table"
import { PaginationState, createColumnHelper } from "@tanstack/react-table"
import { cookies, getCookieValue } from "@/lib/cookies"
import { Badge } from "../ui/badge"

import { ExternalLink } from "lucide-react"
import { Button } from "../ui/button"
import Link from "next/link"
import { UploadDirectoryDemo } from "../dialogs/upload-directory-dialog"

const ch = createColumnHelper<Directory>();

const columns = [
	ch.display({
		id: 'actions',
		cell: props => {
			return (
				<Button variant="link" size="icon">
					<Link href={`/directories/${props.row.original.directory_id}/entries`}>
						<ExternalLink className="h-4 w-4" />
					</Link>
				</Button>
			)
		},
		header: props => {
			return <UploadDirectoryDemo />
		}
	}),
	ch.accessor('directory_name', {
		header: 'File Name'
	}),
	ch.accessor('directory_status', {
		cell: (props) => <Badge>{props.getValue()}</Badge>,
		header: 'Status'
	}),
	ch.accessor('directory_created_at', {
		cell: (props) => new Date(String(props.getValue())).toLocaleTimeString(),
		header: 'Uploaded At'
	})
]

async function getData(props: PaginationState) {
	const res = await client.GET("/directories", {
		params: {
			query: { limit: props.pageSize, offset: props.pageIndex + 1 },
			cookie: { userId: getCookieValue("userId") || "" }
		},
		headers: { Cookie: cookies() },
	})


	if (res.error) {
		throw new Error('Oops something went wrong')
	}

	return res.data.directories
}

export const DirectoriesDataTable = ({ data, rowCount }: { data: Directories, rowCount: number }) => {
	return <DataTable queryKey={'directories'}
		defaultPageSize={100} pageSizes={[25, 50, 100, 250, 500]}
		data={data} columns={columns}
		rowCount={rowCount}
		queryFn={(state) => () => getData(state)} />
}
