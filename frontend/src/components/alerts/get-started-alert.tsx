import { Sparkles } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "../ui/alert";

export default function GetStartedAlert() {
	return (
		<Alert className="bg-stone-200 bg-opacity-45">
			<Sparkles className="w-4 h-4" />
			<AlertTitle>Get Started</AlertTitle>
			<AlertDescription>
				Upload a new directory from the <a href="/" className="text-blue-600">Directories page</a>
			</AlertDescription>
		</Alert>
	)
}


