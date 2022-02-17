package tf6muxserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// schemaType returns the Type for a Schema.
//
// This function should be migrated to a (*tfprotov6.Schema).Type() method
// in terraform-plugin-go.
func schemaType(schema *tfprotov6.Schema) tftypes.Type {
	if schema == nil {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{},
		}
	}

	return schemaBlockType(schema.Block)
}

// schemaAttributeType returns the Type for a SchemaAttribute.
//
// This function should be migrated to a (*tfprotov6.SchemaAttribute).Type()
// method in terraform-plugin-go.
func schemaAttributeType(attribute *tfprotov6.SchemaAttribute) tftypes.Type {
	if attribute == nil {
		return nil
	}

	if attribute.NestedType != nil {
		return schemaObjectType(attribute.NestedType)
	}

	return attribute.Type
}

// schemaBlockType returns the Type for a SchemaBlock.
//
// This function should be migrated to a (*tfprotov6.SchemaBlock).Type()
// method in terraform-plugin-go.
func schemaBlockType(block *tfprotov6.SchemaBlock) tftypes.Type {
	if block == nil {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{},
		}
	}

	attributeTypes := map[string]tftypes.Type{}

	for _, attribute := range block.Attributes {
		if attribute == nil {
			continue
		}

		attributeType := schemaAttributeType(attribute)

		if attributeType == nil {
			continue
		}

		attributeTypes[attribute.Name] = attributeType
	}

	for _, block := range block.BlockTypes {
		if block == nil {
			continue
		}

		blockType := schemaNestedBlockType(block)

		if blockType == nil {
			continue
		}

		attributeTypes[block.TypeName] = blockType
	}

	return tftypes.Object{
		AttributeTypes: attributeTypes,
	}
}

// schemaNestedBlockType returns the Type for a SchemaNestedBlock.
//
// This function should be migrated to a (*tfprotov6.SchemaNestedBlock).Type()
// method in terraform-plugin-go.
func schemaNestedBlockType(nestedBlock *tfprotov6.SchemaNestedBlock) tftypes.Type {
	if nestedBlock == nil {
		return nil
	}

	switch nestedBlock.Nesting {
	case tfprotov6.SchemaNestedBlockNestingModeGroup:
		return schemaBlockType(nestedBlock.Block)
	case tfprotov6.SchemaNestedBlockNestingModeList:
		return tftypes.List{
			ElementType: schemaBlockType(nestedBlock.Block),
		}
	case tfprotov6.SchemaNestedBlockNestingModeMap:
		return tftypes.Map{
			ElementType: schemaBlockType(nestedBlock.Block),
		}
	case tfprotov6.SchemaNestedBlockNestingModeSet:
		return tftypes.Set{
			ElementType: schemaBlockType(nestedBlock.Block),
		}
	case tfprotov6.SchemaNestedBlockNestingModeSingle:
		return schemaBlockType(nestedBlock.Block)
	default:
		return nil
	}
}

// schemaObjectType returns the Type for a SchemaObject.
//
// This function should be migrated to a (*tfprotov6.SchemaObject).Type()
// method in terraform-plugin-go.
func schemaObjectType(object *tfprotov6.SchemaObject) tftypes.Type {
	if object == nil {
		return nil
	}

	attributeTypes := map[string]tftypes.Type{}

	for _, attribute := range object.Attributes {
		if attribute == nil {
			continue
		}

		attributeType := schemaAttributeType(attribute)

		if attributeType == nil {
			continue
		}

		attributeTypes[attribute.Name] = attributeType
	}

	objectType := tftypes.Object{
		AttributeTypes: attributeTypes,
	}

	switch object.Nesting {
	case tfprotov6.SchemaObjectNestingModeList:
		return tftypes.List{
			ElementType: objectType,
		}
	case tfprotov6.SchemaObjectNestingModeMap:
		return tftypes.Map{
			ElementType: objectType,
		}
	case tfprotov6.SchemaObjectNestingModeSet:
		return tftypes.Set{
			ElementType: objectType,
		}
	case tfprotov6.SchemaObjectNestingModeSingle:
		return objectType
	default:
		return nil
	}
}
