import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table"

export function CsvPreviewTable({ rows, excludeFirstRow }: { rows: Array<Array<string>> | undefined, excludeFirstRow: boolean }) {

	if (rows == undefined || rows.length == 0) {
		return (
			<p>Upload a file to see a preview</p>
		)
	}
	return (
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead>First Name</TableHead>
					<TableHead>Last Name</TableHead>
					<TableHead>Email Name</TableHead>
				</TableRow>

			</TableHeader>
			<TableBody>
				{rows.map((row, index) => {
					return (
						<TableRow key={row.join('|')} className={index == 0 && excludeFirstRow ? "line-through bg-gray-200 opacity-55" : ""}>
							<TableCell>{row[0]}</TableCell>
							<TableCell>{row[1]}</TableCell>
							<TableCell>{row[2]}</TableCell>
						</TableRow>
					)
				})}
			</TableBody>
		</Table>
	)
}
