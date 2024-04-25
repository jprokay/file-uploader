"use client"
import { Directories, Directory, client } from "@/lib/contacts-api"
import { DataTable } from "./data-table"
import { ColumnDef, PaginationState, createColumnHelper } from "@tanstack/react-table"
import { cookies, getCookieValue } from "@/lib/cookies"
import { Badge } from "../ui/badge"

import { ExternalLink } from "lucide-react"
import { Button } from "../ui/button"
import Link from "next/link"
import { UploadDirectoryDemo } from "../dialogs/upload-directory-dialog"

const ch = createColumnHelper<Directory>();

const columns: ColumnDef<Directory>[] = [
	{
		accessorKey: 'directory_id',
		id: 'actions',
		cell: props => (
			<Button variant="link" size="icon">
				<Link href={`/directories/${props.getValue()}/entries`}>
					<ExternalLink className="h-4 w-4" />
				</Link>
			</Button>
		),
		header: () => {
			return <UploadDirectoryDemo />
		}
	},
	{
		accessorKey: 'directory_name',
		header: 'File Name',
	},

	{
		accessorKey: 'directory_status',
		cell: (props) => <Badge>{String(props.getValue())}</Badge>,
		header: 'Status'
	},
	{
		accessorKey: 'directory_created_at',
		cell: (props) => new Date(String(props.getValue())).toLocaleTimeString(),
		header: 'Uploaded At'
	}
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
		data={data}
		columns={columns}
		rowCount={rowCount}
		queryFn={(state) => () => getData(state)} />
}
