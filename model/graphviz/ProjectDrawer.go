package graphviz

import (
	"strings"
	model "../core"
)

// EXTEND PROJECT

type ProjectDrawer struct {
	//inherit
	originalProject *model.Project
}

// Decorate Draw function
func (ProjectDrawer *ProjectDrawer) DrawComplete() string {
	var result string
	result = "digraph { graph [] \n"
	// Nodes
	for key, componentList := range ProjectDrawer.originalProject.GetComponentsByGroup() {
		drawingResultInGroup := ""
		for _, component := range componentList {
			drawer := ComponentDrawer{originalComponent: component}
			drawingResultInGroup = drawingResultInGroup + drawer.Draw()
		}

		if key != model.NOGROUP {
			result += "subgraph \"cluster_" + key + "\" { label=\"" + key + "\"; \n " + drawingResultInGroup + "\n}\n"
		} else {
			result += drawingResultInGroup
		}
	}
	// Paths
	for _, component := range ProjectDrawer.originalProject.Components {
		result = result + drawComponentOutgoingRelations(component)
	}
	result = result + "}"
	return result
}



// Decorate Draw function - Draws only a component with its direct dependencies
func (ProjectDrawer *ProjectDrawer) DrawComponent(Component *model.Component) string {
	var result string
	result = "digraph { graph [] \n"
	drawer := ComponentDrawer{originalComponent: Component}
	result = result + drawer.Draw()
	result = result + drawComponentOutgoingRelations(Component)
	allRelatedComponents,_  := Component.GetAllRelatedComponents(ProjectDrawer.originalProject)

	for _, relatedComponent := range allRelatedComponents {
		drawer := ComponentDrawer{originalComponent: &relatedComponent}
		result = result + drawer.Draw()
	}
	// Draw infrastructure :-)
	for _, infrastructureDependency := range Component.InfrastructureDependencies {
		result = result + "\n"+infrastructureDependency.Type+"[shape=box, color=deeppink]"
		result = result + "\n"+infrastructureDependency.Type+"->"+Component.Name+"[color=deeppink,arrowhead=none]"
	}
	result = result + "\n}"
	return result
}


func drawComponentOutgoingRelations(Component *model.Component) string {
	result := ""
	// Relation from components
	for _, dependency := range Component.Dependencies {
		result += "\""+Component.Name + "\" ->" + getGraphVizReference(dependency)+getEdgeLayoutFromDependency(dependency,Component.Display)+"\n"
	}
	// Relation from components/interfaces
	for _, providedInterface := range Component.ProvidedServices {
		for _, dependency := range providedInterface.Dependencies {
			result += "\""+Component.Name+"\":\""+providedInterface.Name+"\"->"+getGraphVizReference(dependency)+getEdgeLayoutFromDependency(dependency,Component.Display)+"\n"
		}
	}
	return result
}

//Get GRaphviz style reference
func getGraphVizReference(Dependency model.Dependency) string {
	if strings.Contains(Dependency.Reference, ".") {
		splitted := strings.Split(Dependency.Reference, ".")
		return "\"" + splitted[0] + "\":\"" + splitted[1] + "\""
	}
	return "\"" + Dependency.Reference + "\""
}

func getEdgeLayoutFromDependency(dependency model.Dependency, display model.ComponentDisplaySettings) string {

	edgeLayout := "["
	if display.BorderColor != "" {
		edgeLayout += "color=\"" + display.BorderColor + "\""
	} else {
		edgeLayout += "color=\"#333333\""
	}
	edgeLayout += ", fontsize=\"10\", fontcolor=\"#555555\" "
	if dependency.Relationship == "acl" {
		edgeLayout += ", dir=both, arrowtail=\"box\", taillabel=<<font color=\"red\" ><b>acl</b></font>>"
	} else {
		edgeLayout += ", label=\"" + dependency.Relationship + "\""
	}
	if dependency.Relationship == "customer-supplier" || dependency.Relationship == "conformist" {
		edgeLayout += "weight=3, style=\"bold\""
	} else if dependency.IsBrowserBased {
		edgeLayout += "style=\"dashed\""
	}
	if dependency.IsSameLevel {
		edgeLayout += ", constraint=false"
	}
	return edgeLayout + "]"
}

// Factory
func CreateProjectDrawer(Project *model.Project) *ProjectDrawer {
	var Drawer ProjectDrawer
	Drawer.originalProject = Project
	return &Drawer
}
