import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Github, Linkedin, Mail } from "lucide-react"

export default function AboutPage() {
  return (
    <div className="flex flex-col gap-6 p-6">
      <div>
        <h1 className="text-3xl font-bold">About This Project</h1>
        <p className="text-muted-foreground mt-2">
          A cloud-based image storage service
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Project Overview</CardTitle>
          <CardDescription>Modern cloud storage solution for images</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground leading-relaxed">
            This project is a full-stack cloud-based image storage service that combines the power 
            of AWS infrastructure with modern web technologies. Built to handle scalable storage 
            needs with a clean and intuitive user interface.
          </p>
          
          <div>
            <h4 className="font-semibold mb-2">Backend Technologies</h4>
            <ul className="text-sm text-muted-foreground space-y-1 ml-4">
              <li>• Go with Gin framework for high-performance API</li>
              <li>• AWS S3 for reliable image storage</li>
              <li>• AWS DynamoDB for metadata management</li>
              <li>• AWS Lambda/EC2 for serverless compute</li>
            </ul>
          </div>

          <div>
            <h4 className="font-semibold mb-2">Frontend Technologies</h4>
            <ul className="text-sm text-muted-foreground space-y-1 ml-4">
              <li>• Next.js 16 for server-side rendering</li>
              <li>• shadcn/ui for beautiful components</li>
              <li>• Tailwind CSS for responsive design</li>
              <li>• TypeScript for type safety</li>
            </ul>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Get in Touch</CardTitle>
          <CardDescription>Connect with me on social platforms</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-3">
          <Button variant="outline" size="sm" asChild>
            <a href="https://github.com/berkkaradalan" target="_blank" rel="noopener noreferrer">
              <Github className="mr-2 h-4 w-4" />
              GitHub
            </a>
          </Button>
          
          <Button variant="outline" size="sm" asChild>
            <a href="https://linkedin.com/in/berkkaradalan" target="_blank" rel="noopener noreferrer">
              <Linkedin className="mr-2 h-4 w-4" />
              LinkedIn
            </a>
          </Button>
          
          <Button variant="outline" size="sm" asChild>
            <a href="mailto:berkkaradalan@gmail.com">
              <Mail className="mr-2 h-4 w-4" />
              Email
            </a>
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}